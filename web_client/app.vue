<template>
  <div>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script lang="ts">
useHead({
  title: 'オンライン作業部屋',
  meta: [
    { charset: 'utf-8' },
    { name: 'viewport', content: 'width=device-width, initial-scale=1' },
    {
      hid: 'description',
      name: 'description',
      content:
        'オンラインの作業部屋で集中しよう。Googleアカウントでログインして、各科目・作業内容の自習室を選んで入室しましょう。サービスは全て無料です。',
    },
    {
      hid: 'keywords',
      name: 'keywords',
      content: 'オンライン作業部屋,オンライン自習室,オンライン勉強部屋,study with me,作業耐久',
    },
  ],
  script: [],
  link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
})
</script>
