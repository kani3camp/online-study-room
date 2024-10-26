<template>
  <v-dialog v-model="ifShow" width="500">
    <v-card class="mx-auto" outlined :loading="loading">
      <v-card-title> {{ cardTitle }}</v-card-title>

      <v-card-actions box-sizing>
        <v-spacer />
        <v-btn
          v-if="ifAcceptNeeded"
          :disabled="loading"
          text="true"
          color="primary"
          @click="handleAccept"
        >
          {{ acceptOptionString }}
        </v-btn>
        <v-btn :disabled="loading" pr-0 text="true" @click="handleCancel">
          {{ cancelOptionString }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, computed } from 'vue'

const props = defineProps({
  ifShowDialog: {
    type: Boolean,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  cardTitle: {
    type: String,
    required: true,
  },
  acceptOptionString: {
    type: String,
    default: '',
  },
  cancelOptionString: {
    type: String,
    default: 'キャンセル',
  },
  acceptNeeded: {
    type: Boolean,
    required: true,
  },
})

const emit = defineEmits(['accept', 'cancel'])

const ifShow = computed({
  get: () => props.ifShowDialog,
  set: (value: boolean) => {
    // 親コンポーネントに変更を通知する場合はemitを使用
    // ここでは読み取り専用として設定
  },
})

const handleAccept = () => {
  emit('accept')
}
const handleCancel = () => {
  emit('cancel')
}
const ifAcceptNeeded = computed(() => props.acceptNeeded)
</script>

<style scoped></style>
