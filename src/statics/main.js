
var View = Vue.extend({
  template: '#view',
  data: function () {
    return {
        err: '',
        isDisabled: false,
        token: getCookie("token"),
        content: '',
        text: '',
        decrypted: false,
        err2: ''
        // entryId: this.$route.query.entry_id
    }
  },
  created: function () {
    this.fetchData();
  },
  methods: {
    fetchData: function() {
      var json = {"text":this.$route.params.entry_id, "key":""}
      
      this.$http.post('/api/entry', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        this.content = response.body;
        this.err2 = this.content.type
        this.text = this.content.EncryptedText
      })
    },
    decrypt: function() {
      key = this.$refs[this.content.Id].value;
      var json = {"text":this.content.EncryptedText, "key":key}
      this.$http.post('/api/decrypt', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        this.text = response.body;
      })
    }
  }
})

var Home = Vue.extend({
  template: '#home-list',
  data: function () {
    return {searchKey: '', data: [], key:'', text:''};
  },
  created: function () {
    this.fetchData();
  },
  methods : {
      fetchData: function () {

        this.$http.get('/api/list', {headers: {'Content-Type':'application/json', 'Authorization': 'Bearer '+getCookie("token")}}).then(response => {

        res = response.body;
        if (res == "\"this token is not authorized for this content\"") {

          document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
          router.push('/login');  

        } else {
          k = this.searchKey;
          this.data = res.filter(function (product) {
            return product.Day.indexOf(k) !== -1
          })
        }
        })
      },
  },
  computed : {
    data2: function () {
      var self = this;
      return self.data.filter(function (product) {
        return product.Day.indexOf(self.searchKey) !== -1
      })
    }
  }
});





var Register = Vue.extend({
  template: '#register',
  data: function () {
    return {
        user: {username: '', password: ''},
        err: '',
        isDisabled: false,
        token: getCookie("token"),

    }
  },
  methods: {
    createUser: function() {
      var user = this.user;
      var json = {"username": user.username, "password":user.password}
      var res = "";
      this.$http.post('/api/register', json, {
        headers: {
        'Content-Type': 'application/json'
        }
      }).then(response => {
        res = response.body;

        if (res=="\"success\""){
          this.err = "success"
          this.isDisabled = true
          setTimeout(function(){ 
            router.push('/login');  
          }, 3000);
          
        } else if (res=="\"already registered username\""){
          this.err = "alreadyuser"

        }
      })
      
      
    }
  }
});

var Login = Vue.extend({
  template: '#login',
  data: function () {
    return {
        user: {username: '', password: ''},
        err: '',
        isDisabled: false,
        token: getCookie("token"),//getFromLocalStorage("Token")

    }
  },
  methods: {
    logInUser: function() {
      var user = this.user;
      var json = {"username": user.username, "password":user.password}
      var res = "";
      this.$http.post('/api/login', json, {
        headers: {
        'Content-Type': 'application/json'
        }
      }).then(response => {
        res = response.body;

        if (res=="\"user doesn't exist\""){
          this.err = "doesn't exist";
        } else if (res=="\"password is not true\""){
          this.err = "incorrectpass";

        } else if (res.substring(0, 10)=="\"Error while".substring(0, 10) || res.substring(0, 10) == "\"everything is ".substring(0, 10)){
          this.err = "othererr";

        } else {
          this.err = "success";
          
          document.cookie = "token=" + res.substring(1, res.length-1);
          //saveToLocalStorage("Token",res.substring(1, res.length-1));
          
          this.isDisabled= false;
          setTimeout(function(){ 
            router.push('/');  
          }, 3000);
        }
      })
      
      
    }
  }
});


var Post = Vue.extend({
  template: '#post',
  data: function () {
    return {
        post: {text: '', key: ''},
        err: '',
        isDisabled: false,
        token: getCookie("token"),
    }
  },
  methods: {
    postEntry: function() {
      var post = this.post;
      var json = {"text": post.text, "key":post.key}
      var res = "";
      this.$http.post('/api/post', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        res = response.body;

        if (res=="\"success\""){
          this.err = "success"
          this.isDisabled = true
          setTimeout(function(){ 
            router.push('/');  
          }, 3000);
          
        } else {
          this.err = "error"

        }
      })
      
      
    }
  }
});

var saveToLocalStorage = function(k,v){
  localStorage.setItem(k, v);
}

var getFromLocalStorage = function(k){
  res = localStorage.getItem(k);
  if (!res){
    return ""
  }
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return '';
}

// var isTokenValid = function(){
//   token = "Bearer " + getFromLocalStorage("Token");
//   console.log(token);
//   var xhttp = new XMLHttpRequest();
//   xhttp.open("GET", "/token", true);
//   xhttp.setRequestHeader("Content-Type", "application/json");
//   xhttp.setRequestHeader("Authorization", token);
//   xhttp.send();
//   xhttp.onload = function () {
//     Token = this.responseText;
//   };
// }


const router = new VueRouter({
  routes: [
    {path: '/', component: Home, name:'home'},
    {path: '/post', component: Post, name:'post'},
    {path: '/register', component: Register, name: 'register'},
    {path: '/login', component: Login, name: 'login'},
    {path: '/view/:entry_id', component: View, name: 'view'},
  ]
});

new Vue({
  el: '#app',
  router: router,
  template: '<router-view></router-view>'
});

