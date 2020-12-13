<template>
  <v-dialog
    v-model="ifShow"
    width="500"
  >
    <v-card
      class="mx-auto"
      outlined
      :loading="loading"
    >
      <v-card-title> {{ cardTitle }}</v-card-title>

      <v-card-actions box-sizing>
        <v-spacer />
        <v-btn
          v-if="ifAcceptNeeded"
          :disabled="loading"
          text
          color="primary"
          @click="ifAccept=true"
        >
          {{ acceptOptionString }}
        </v-btn>
        <v-btn
          :disabled="loading"
          pr-0
          text
          @click="ifCancel=true"
        >
          {{ cancelOptionString }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'Dialog',
  props: {
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
  },
  computed: {
    ifShow: {
      get() {
        return this.ifShowDialog
      },
      set() {},
    },
    ifAccept: {
      get() {
        return false
      },
      set(value) {
        if (value) {
          this.$emit('accept')
        }
      },
    },
    ifCancel: {
      get() {
        return false
      },
      set(value) {
        if (value) {
          this.$emit('cancel')
        }
      },
    },
    ifAcceptNeeded: {
      get() {
        return this.acceptNeeded
      },
      set() {},
    },
  },
}
</script>

<style scoped></style>
