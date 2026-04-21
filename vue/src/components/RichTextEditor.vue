<template>
  <div class="editor-wrapper">
    <div ref="editorRef"></div>
    <input ref="fileInputRef" type="file" accept="image/*" class="hidden-input" @change="onFileChange" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import Quill from 'quill'
import 'quill/dist/quill.snow.css'
import request from '@/utils/request'
import ImageResize from 'quill-image-resize-module'

// 注册图片缩放模块
Quill.register('modules/imageResize', ImageResize)

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = ref(null)
const fileInputRef = ref(null)

let quill = null
let lastValue = ''

const openFileDialog = () => {
  if (!fileInputRef.value) return
  fileInputRef.value.value = ''
  fileInputRef.value.click()
}

const uploadImage = async (file) => {
  const formData = new FormData()
  formData.append('contentImage', file)
  const res = await request.post('/api/article/contentImage', formData)
  return res.data
}

const onFileChange = async (e) => {
  const file = e.target.files && e.target.files[0]
  if (!file || !quill) return
  const range = quill.getSelection(true)
  try {
    const url = await uploadImage(file)
    quill.insertEmbed(range ? range.index : 0, 'image', url, 'user')
    quill.setSelection((range ? range.index : 0) + 1, 0, 'silent')
  } catch (err) {
    console.error(err)
  }
}

onMounted(() => {
  const toolbarOptions = [
    [{ header: [1, 2, 3, 4, 5, 6, false] }],
    ['bold', 'italic', 'underline', 'strike'],
    ['blockquote', 'code-block'],
    [{ list: 'ordered' }, { list: 'bullet' }],
    [{ color: [] }, { background: [] }],
    ['link', 'image'],
    ['clean']
  ]

  quill = new Quill(editorRef.value, {
    theme: 'snow',
    placeholder: '在这里开始编写你的文章...',
    modules: {
      toolbar: {
        container: toolbarOptions,
        handlers: {
          image: openFileDialog
        }
      },
      imageResize: {
        displayStyles: {
          backgroundColor: 'black',
          border: 'none',
          color: 'white'
        },
        modules: ['Resize', 'DisplaySize', 'Toolbar']
      }
    }
  })

  quill.on('text-change', () => {
    const html = quill.root.innerHTML
    lastValue = html
    emit('update:modelValue', html)
  })

  if (props.modelValue) {
    quill.root.innerHTML = props.modelValue
    lastValue = props.modelValue
  }
})

watch(
    () => props.modelValue,
    (val) => {
      if (!quill) return
      if (val === lastValue) return
      quill.root.innerHTML = val || ''
      lastValue = val || ''
    }
)

onBeforeUnmount(() => {
  quill = null
})
</script>

<style scoped>
.editor-wrapper {
  min-height: 400px;
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
}

.editor-wrapper {
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  background-color: #f8fafc;
}

.hidden-input {
  display: none;
}

/* 图片缩放模块样式 */
.editor-wrapper {
  max-width: 100%;
  cursor: pointer;
}

.editor-wrapper {
  display: inline-block;
  position: relative;
}

.editor-wrapper {
  opacity: 0.5;
}

.editor-wrapper {
  position: absolute;
  width: 10px;
  height: 10px;
  background: #409eff;
  border: 2px solid #fff;
  border-radius: 50%;
  z-index: 100;
}

.editor-wrapper {
  right: -5px;
  bottom: -5px;
  cursor: se-resize;
}

.editor-wrapper {
  left: -5px;
  bottom: -5px;
  cursor: sw-resize;
}

.editor-wrapper {
  right: -5px;
  top: -5px;
  cursor: ne-resize;
}

.editor-wrapper {
  left: -5px;
  top: -5px;
  cursor: nw-resize;
}
</style>