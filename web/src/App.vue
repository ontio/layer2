<template>
  <div id="app">
      <el-row type="flex" justify="center">
          <el-col :xs="24" :sm="24" :lg="8">
            <router-view/>
          </el-col>
      </el-row>
  </div>
</template>

<style lang="scss">
@media screen and (min-width: 420px) {
    #app {
        width: 100%;
    }
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
p {
    margin: 0;
    padding: 0;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
