<template>
  <pev2 v-if="planSource" :plan-source="planSource" :plan-query="planQuery" />
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import pev2 from "pev2";
import router from '../router';

import { planData } from '../App.vue';

@Component({
  components: {
    pev2: pev2,
  },
})
export default class App extends Vue {
  private planSource: any | any[] = planData[0];
  private planQuery: string = planData[1];

  private beforeMount() {
    if (!this.planSource) {
      router.push({ name: 'home' });
    }
  }
}
</script>
