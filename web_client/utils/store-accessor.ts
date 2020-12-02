/* eslint-disable import/no-mutable-exports */
import { Store } from 'vuex'
import { getModule } from 'vuex-module-decorators'
import User from '~/store/user'

let UserStore: User
function initializeStores(store: Store<any>): void {
  UserStore = getModule(User, store)
}

export { initializeStores, UserStore }
