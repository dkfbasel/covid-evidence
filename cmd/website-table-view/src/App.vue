<template>
  <div style="width:100%;height:600px;overflow:hidden" v-if="rows && rows.length > 0">
    <ejs-grid ref="grid" :dataSource="rows"
      height="100%" width="100%"
      :show-column-chooser="true"
      :toolbar="settings.toolbar"
      :toolbar-click="handleToolbarClick"
      :allow-sorting="true"
      :allow-filtering="true"
      :allow-resizing="true"
      :allow-excel-export="true"
      :filter-settings="settings.filter"
      :sort-settings="settings.sorting"
      :allow-text-wrap="false"
      :query-cell-info="showTooltip"
      row-height="40"
      :frozenColumns="1"> 
      <e-columns>
        <e-column v-for="col in columns" 
          :key="col.field" 
          :field="col.field" :header-text="col.label || col.field"
          :filter="col.filter || { type:'CheckBox' }"
          :format="col.format"
          :sort-comparer="col.sorting || null"
          :allow-filtering="col.allowFiltering == false ? false : true"
          :width="col.width"
          :disable-html-encode="col.html ? false : true"
          :clip-mode="settings.clipMode" />
      </e-columns>
    </ejs-grid>
  </div>
</template>

<script>

const updateDate = '20200409';

// import axios from 'axios';
import data from './data/data_20200505_filtered.json';

import { Sort, Filter, Group, Freeze, ColumnChooser, Toolbar, Search, Resize, ExcelExport  } from "@syncfusion/ej2-vue-grids";
import { Tooltip } from "@syncfusion/ej2-popups";

let reviewStatusSortOrder = ['review completed', 'in review', 'manual extraction completed', 'in manual extraction', 'prefilled automatically', '', undefined];

let sortBy = function(sortOrder) {
  return function(a, b) {
    return sortOrder.indexOf(a) - sortOrder.indexOf(b);
  }
}

let is_safari = navigator.userAgent.indexOf("Safari") > -1;
let is_chrome = navigator.userAgent.indexOf('Chrome') > -1;

if ((is_chrome)&&(is_safari)) { is_safari = false; }

export default {
  name: 'DatabaseGrid',
  // provide additional functionality to grid component
  provide: {
    grid: [Sort, Filter, Group, Freeze, ColumnChooser, 
      Toolbar, Search, Resize, ExcelExport]
  },
  data() {
    return {
      settings: {
        filter: {
          ignoreAccent: true,
          type: "CheckBox"
        },
        sorting: {
          columns: [
            { field: 'results_available', direction:'Ascending'},
            { field: 'review_status', direction: 'Ascending'}
          ]
        },
        wrap: {
          wrapMode: 'Content'
        },
        clipMode: 'EllipsisWithTooltip',
        toolbar: [ 'Search'],
      },

      // column information
      columns: [{
        label: "ID",
        field: "cove_id",
        width: 94,
      }, {
        field: "source",
        html: true,
        width: 180
      }, {
        label: "results available",
        field: "results_available",
        sorting: sortBy(['yes', 'unclear', 'no', '', undefined]),
        width: 120
      }, {
        label: "status of review",
        field: "review_status",
        sorting: sortBy(reviewStatusSortOrder),
        width: 180
      }, {
        label: "individual patient data sharing",
        field: "ipd_sharing",
        width: 120
      }, {
        label: "intervention type",
        field: "intervention_type",
        width: 120
      }, {
        label: "intervention name",
        field: "intervention_name",
        filter: { type: 'Menu' },
        width: 250
      }, {
        label: "number of participants entrolled",
        field: "n_enrollment",
        format: 'N0',
        filter: { type: 'Menu', operator:"greaterThan" },
        width: 100
      }, {
        field: "country",
        width: 120
      }, {
        field: "status",
        width: 150
      }, {
        field: "randomized",
        width: 150
      }, {
        label: "arms",
        field: "n_arms",
        format: 'N0',
        width: 110
      }, {
        field: "blinding",
        width: 120
      }, {
        label: "population condition",
        field: "population_condition",
        filter: { type: 'Menu' },
        width: 200
      }, {
        field: "control",
        filter: { type: 'Menu' },
        width: 200
      }, {
        label: "primary outcome measure",
        field: "out_primary_measure",
        filter: { type: 'Menu' },
        width: 200
      }, {
        label: "start date",
        field: "start_date",
        filter: { type: 'Menu' },
        width: 120
      }, {
        label: "end date",
        field: "end_date",
        filter: { type: 'Menu' },
        width: 120
      }, {
        field: "source_id",
        filter: { type: 'Menu' },
        width: 150
      }, {
        field: "title",
        filter: { type: 'Menu' },
        width: 250
      }, {
        field: "abstract",
        filter: { type: 'Menu' },
        width: 250
      }, {
        label: "type of entry",
        field: "entry_type",
        width: 120
      }, {
        field: "url",
        filter: { type: 'Menu'},
        width: 250
      }, ],
      // row data
      rows: [],
    }
  },
  mounted() {

    let dta = data;

    for (let i = 0; i < dta.length; i++) {
      if (!dta[i].source) continue;
      if (!dta[i].url) continue;
      dta[i].source = `<a href="${dta[i].url}" target="_blank">${dta[i].source}</a>`;
    }

    this.rows = Object.freeze(dta);

  },
  methods: {
    showTooltip(args) {

      if (!is_safari) {
        return;
      }

      // make sure that the field is actually available in the data
      if (!args.data[args.column.field]) return;

      let content = args.data[args.column.field].toString();

      if (content.length < 25) return;

      new Tooltip({
        content: content,
        position: 'BottomCenter',
        opensOn: 'Hover Click',
        cssClass: 'coveTooltip'
      }, args.cell);

    },

    handleToolbarClick(args) {
      switch (args.item.text) {
          case 'Excel Export':
              (this.$refs.grid).excelExport({
                fileName: `covid-evidence-${updateDate}.xlsx`
              });
              break;
          case 'CSV Export':
              (this.$refs.grid).csvExport({
                fileName: `covid-evidence-${updateDate}.csv`
              });
              break;
      }
    }
  }
}
</script>

<style lang="stylus">
@import '../node_modules/@syncfusion/ej2-base/styles/material.css';
@import '../node_modules/@syncfusion/ej2-buttons/styles/material.css';
@import '../node_modules/@syncfusion/ej2-calendars/styles/material.css';
@import '../node_modules/@syncfusion/ej2-dropdowns/styles/material.css';
@import '../node_modules/@syncfusion/ej2-inputs/styles/material.css';
@import '../node_modules/@syncfusion/ej2-navigations/styles/material.css';
@import '../node_modules/@syncfusion/ej2-popups/styles/material.css';
@import '../node_modules/@syncfusion/ej2-splitbuttons/styles/material.css';
@import "../node_modules/@syncfusion/ej2-vue-grids/styles/material.css";

$color-red = #D20537;
$color-blue = #150f48;

$color-tooltip = $color-blue;

.e-tooltip-wrap {

  &.e-popup {
    background-color: $color-tooltip;
    padding: 8px;
    border-radius: 1px;
    border-color: $color-tooltip;
    opacity: 1;
  }

  .e-tip-content {
    color: #fff;
    font-size: 12px;
    line-height: 12 * 1.4px;
  }
}

.e-arrow-tip-inner {
  color: $color-tooltip !important;
}

.e-tooltip-wrap .e-arrow-tip-outer.e-tip-bottom {
    border-top-color: $color-tooltip !important;
}

.e-tooltip-wrap .e-arrow-tip-outer.e-tip-top {
    border-bottom-color:$color-tooltip !important;
}

.e-tooltip-wrap .e-arrow-tip-outer.e-tip-left {
    border-right-color: $color-tooltip !important;
}

.e-tooltip-wrap .e-arrow-tip-outer.e-tip-right {
    border-left-color: $color-tooltip !important;
}
</style>
