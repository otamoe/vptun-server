<template>
  <main id="main" class="client-main">
    <RouterView></RouterView>
  </main>
</template>
<style lang="scss">
.client-main {
    min-width: 1280px;
    width: 96%;
    margin: 3rem auto 2% auto;
    h1 {
      text-align: center;
      font-size: 2em;
      margin: .5em 0;
    }

  form {
      background: #fff;
      border-radius: .5em;
  }

  table {
    margin: 0;
      vertical-align: top;
    caption-side: bottom;
    border-collapse: collapse;
    width: 100%;
    thead {
          text-align: inherit;
      background: #f8f8f8;
      vertical-align: bottom;
    }
    tbody {
      vertical-align: bottom;
    }
    tr, th, td {
      padding: 0;
      margin: 0;
          text-align: inherit;

    }
    tr {
      
    }
     th, td {
        border: 1px solid #dfe2e5; 
       padding: .25em .5em;
     }
  }
}
</style>