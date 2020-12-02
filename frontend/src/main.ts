import axios from 'axios';
import Vue from 'vue';
import App, { planData } from './App.vue';
import router from './router';
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import '@fortawesome/fontawesome-free/css/all.css';

Vue.config.productionTip = false;

export function setPlanData (plan: string, query: string) {
  planData[0] = plan;
  planData[1] = query;
  axios
    .post('api/', planData)
    .then(response => {
      router.push({ path: 'plan/' + response.data.id });
    })
};

export function loadPlanData(id: string) {
  return axios
    .get('api/' + id)
    .then(response => {
      planData[0] = response.data[0];
      planData[1] = response.data[1];
    })
}

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app');
