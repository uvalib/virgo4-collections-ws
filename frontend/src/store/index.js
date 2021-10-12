import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
   state: {
      working: false,
      collections: [],
      selectedID: -1,
      fatal: ""
   },
   mutations: {
      setCollections(state, data) {
         state.collections.splice(0, state.collections.length)
         data.forEach(c => state.collections.push(c))
      },
      setFatalError(state, err) {
         state.fatal = err
      },
      setWorking(state, flag) {
         state.working = flag
      }
   },
   actions: {
      getCollections(ctx) {
         ctx.commit("setWorking", true)
         axios.get("/api/collections").then(resp => {
            ctx.commit("setCollections", resp.data)
            ctx.commit("setWorking", false)
         }).catch((e) => {
            ctx.commit("setFatalError", e)
            ctx.commit("setWorking", false)
         })
      }
   },
})
