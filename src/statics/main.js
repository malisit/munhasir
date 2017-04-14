var Register = Vue.extend({
  template: '#register',
  data: function () {
    return {
        user: {username: '', password: ''},
        err: '',
        isDisabled: false,

    }
  },
  methods: {
    createUser: function() {
      var user = this.user;
      var json = {"username": user.username, "password":user.password}
      var res = "";
      this.$http.post('/register', json, {
        headers: {
        'Content-Type': 'application/json'
        }
      }).then(response => {
        res = response.body;

        if (res=="\"success\""){
          this.err = "success"
          this.isDisabled = true
          setTimeout(function(){ 
            router.push('/');  
          }, 3000);
          
        } else if (res=="\"already registered username\""){
          this.err = "alreadyuser"

        }
      })
      
      
    }
  }
});


var router = new VueRouter({
  routes: [
  {path:   '/register', component: Register, name: 'register'},
]});

new Vue({
  el: '#app',
  router: router,
  template: '<router-view></router-view>'
});