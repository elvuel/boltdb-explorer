<script>
  import { onMount } from "svelte"
  import { JSONEditor } from "svelte-jsoneditor"

  import { upload } from "./api/upload"
  import { download } from "./api/download"

  import Modal from "./lib/Modal.svelte"

  let mounted = false
  let items = []

  let fileVar
  let fileVal

  let content = {
    text: undefined,
    json: {},
  }
  let contentIdx = -1

  let modalShow = false
  let modalTitle = ""
  let modalContent = ""
  let modalKind = ""

  function addItem(data) {
    if (items.find((item) => item.filename === data.filename) === undefined) {
      items = [...items, data]
    }
  }

  function dropItem(data) {
    items = [...items.filter((item) => item.filename !== data.filename)]
  }

  onMount(() => {
    const data = window.localStorage.getItem("items")
    if (data) {
      items = JSON.parse(data)
    }
    mounted = true
  })

  $: items,
    mounted ? window.localStorage.setItem("items", JSON.stringify(items)) : null
</script>

<div class="flex bg-gray-700 text-white">
  <span class="pl-1 font-bold">BoltDB Explorer</span>
</div>
<div class="flex">
  <div class="basis-9/12 h-screen">
    <JSONEditor
      bind:content
      onChange={(updatedContent) => {
        // console.log(updatedContent)
        items[contentIdx].data = updatedContent
      }}
    />
  </div>
  <div class="basis-3/12 h-screen">
    <div class="flex flex-col">
      <div class="flex justify-center">
        <div class="mb-3 w-96">
          <label
            for="formFile"
            class="nline-block mb-2 text-green-500 font-bold">Load BoltDB</label
          >
          <div class="flex">
            <input
              bind:files={fileVar}
              bind:this={fileVal}
              class="
            block
            w-full
            px-3
            py-1
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            shadow-sm
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
              type="file"
              id="formFile"
            />
            <button
              class="px-4 py-2 font-semibold text-sm bg-sky-500 text-white shadow-sm rounded"
              on:click={(e) => {
                e.preventDefault()
                if (fileVal.value === "") {
                  return
                }
                upload(
                  fileVar[0],
                  (data) => {
                    if (data.error) {
                      modalKind = "error"
                      modalTitle = "Upload Failed"
                      modalContent = data.error
                      modalShow = true
                      return
                    }
                    fileVal.value = ""
                    addItem(data)
                  },
                  (err) => {
                    modalKind = "error"
                    modalTitle = "Upload Error"
                    modalContent = err.message
                    modalShow = true
                  }
                )
              }}
            >
              Upload
            </button>
          </div>
        </div>
      </div>
      {#each items as item, i}
        <div class="flex justify-center shadow-sm mb-1">
          <div class="mb-3 w-96 space-x-1">
            <button
              class="nline-blocktext-green-500 font-bold"
              on:click={(e) => {
                e.preventDefault()
                content.text = undefined
                content.json = item.data
                content = { ...content }
                contentIdx = i
              }}
            >
              <span class="text-blue-500">Explore</span
              >&lt;&lt;{item.filename}&gt;&gt;
            </button>
            <button
              class="rounded-md bg-green-500 text-white w-5"
              on:click={(e) => {
                e.preventDefault()
                download(
                  item.data,
                  async (response) => {
                    if (response.status !== 200) {
                      const msg = await response.text()
                      throw new Error(msg)
                    }
                    return response.blob()
                  },
                  (blob) => {
                    const url = window.URL.createObjectURL(blob)
                    let a = document.createElement("a")
                    a.href = url
                    a.download = item.filename
                    document.body.appendChild(a)
                    a.click()
                    a.remove()
                  },
                  (err) => {
                    modalKind = "error"
                    modalTitle = "Download Failed"
                    modalContent = err.message
                    modalShow = true
                  }
                )
              }}
              >&darr;
            </button>
            <button
              class="rounded-md bg-red-500 text-white w-5"
              on:click={(e) => {
                e.preventDefault()
                dropItem(item)
                if (contentIdx === i) {
                  content.text = undefined
                  content.json = {}
                  content = { ...content }
                  contentIdx = -1
                }
              }}>x</button
            >
          </div>
        </div>
      {/each}
    </div>
  </div>

  <Modal
    bind:show={modalShow}
    bind:title={modalTitle}
    bind:content={modalContent}
    bind:kind={modalKind}
  />
</div>
