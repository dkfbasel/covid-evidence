import Vue from 'vue'

// use custom components from buefy (buefy.org)
// import { Table, Input} from 'buefy'
// import 'buefy/dist/buefy.css'

import { GridPlugin } from '@syncfusion/ej2-vue-grids'
import '@syncfusion/ej2-base/styles/material.css'
import '@syncfusion/ej2-buttons/styles/material.css'
import '@syncfusion/ej2-calendars/styles/material.css'
import '@syncfusion/ej2-dropdowns/styles/material.css'
import '@syncfusion/ej2-inputs/styles/material.css'
import '@syncfusion/ej2-navigations/styles/material.css'
import '@syncfusion/ej2-popups/styles/material.css'
import '@syncfusion/ej2-splitbuttons/styles/material.css'
import '@syncfusion/ej2-vue-grids/styles/material.css'

Vue.use(GridPlugin)

// Vue.use(Table)
// Vue.use(Input)

import DatabaseGrid from './App.vue'

Vue.config.productionTip = false

new Vue({
  render: h => h(DatabaseGrid),
}).$mount('#database-grid')
