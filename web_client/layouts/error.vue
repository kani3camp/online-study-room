<template>
  <div>
    <h1 v-if="error.statusCode === 404">ページが見つかりませんでした。</h1>
    <h1 v-else>エラーが発生しました。</h1>
    <a href="https://online-study-space.web.app">トップページへ</a>
  </div>
</template>

<script>
export default {
  layout: 'error',
  props: ['error'], // エラーページ用のカスタムレイアウトを指定できます
}
</script>
