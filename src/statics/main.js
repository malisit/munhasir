var changePassword = Vue.extend({
  template: '#change-password',
  data: function(){
    return {
      token: getCookie("token"),
      oldPassword1: '',
      oldPassword2: '',
      oldPasswordTemp: '',
      newPassword: '',
      status: ''
    }
  },
  methods: {
    post: function(){
      var json = {"text":this.oldPassword1, "key":this.newPassword}
      
      this.$http.post('/api/user/password', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        this.status = response.body.substring(1,response.body.length-1);
        this.oldPasswordTemp = this.oldPassword1;
        if (this.status == 'success') {

          document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
          setTimeout(function(){ 
            router.push('/login');  
          }, 3000);
        }
      })
    }
  }
})


var summernoteComponent = {
    replace: true,
    inherit: false,
    template: "<textarea class='form-control' :name='name'></textarea>",
    props: {
        model: {
            required: true,
            twoWay: true
        },
        language: {
            type: String,
            required: false,
            default: "en-US"
        },
        height: {
            type: Number,
            required: false,
            default: 160
        },
        minHeight: {
            type: Number,
            required: false,
            default: 160
        },
        maxHeight: {
            type: Number,
            required: false,
            default: 800
        },
        name: {
            type: String,
            required: false,
            default: ""
        },
        toolbar: {
            type: Array,
            required: false,
            default: function() {
                return [
                    ["font", ["font","bold", "italic", "underline", "clear"]],
                    ["style", ["fontname","strikethrough", "superscript", "subscript"]],
                    ["fontsize", ["fontsize"]],
                    ["para", ["ul", "ol", "paragraph"]],
                    ["insert", ["link", "hr"]],
                    ['height', ['height','codeview']]
                ];
            }
        }
    },
    created: function() {
        this.isChanging = false;
        this.control = null;
    },
    mounted: function() {
        //  initialize the summernote
        if (this.minHeight > this.height) {
            this.minHeight = this.height;
        }
        if (this.maxHeight < this.height) {
            this.maxHeight = this.height;
        }
        var me = this;
        this.control = $(this.$el);
        this.control.summernote({
            lang: this.language,
            height: this.height,
            minHeight: this.minHeight,
            maxHeight: this.maxHeight,
            toolbar: this.toolbar,
            callbacks: {
                onInit: function() {
                    me.control.summernote("code", me.model);
                },
                onChange: function() {
                    if (!me.isChanging) {
                        me.isChanging = true;
                        var code = me.control.summernote("code");
                        me.model = (code === null || code.length === 0 ? null : code);
                        me.$nextTick(function() {
                            me.isChanging = false;
                        });
                    }
                    me.$parent.text = code

                }
            }
        })
    },
    watch: {
        'model': function(val) {
            if (!this.isChanging) {
                
                this.isChanging = true;
                var code = (val === null ? "" : val);
                this.control.summernote("code", code);
                this.isChanging = false;
            }
        }
    },
}

// var View = Vue.extend({
//   template: '#view',
//   data: function () {
//     return {
//         err: '',
//         isDisabled: false,
//         token: getCookie("token"),
//         content: '',
//         text: '',
//         decrypted: false,
//         err2: '',
//         pass: ''
//         // entryId: this.$route.query.entry_id
//     }
//   },
//   created: function () {
//     this.fetchData();
//   },
//   methods: {
//     fetchData: function() {
//       var json = {"text":this.$route.params.entry_id, "key":""}
      
//       this.$http.post('/api/entry/view', json, {
//         headers: {
//         'Content-Type': 'application/json',
//         'Authorization': 'Bearer '+getCookie("token")
//         }
//       }).then(response => {
//         this.content = response.body;
//         this.err2 = this.content.type
//         this.text = this.content.EncryptedText
//       })
//     },
//     decrypt: function() {
//       key = this.$refs[this.content.Id].value;
//       var json = {"text":this.content.EncryptedText, "key":key}
//       this.$http.post('/api/entry/decrypt', json, {
//         headers: {
//         'Content-Type': 'application/json',
//         'Authorization': 'Bearer '+getCookie("token")
//         }
//       }).then(response => {
//         this.text = response.body.substring(1,response.body.length-1);
//       })
//     },
    
//   }
// })

var DeleteEntry = Vue.extend({
  template: '#del-entry',
  data: function() {
    return {
      pass: '',
      token: getCookie("token"),
      status: '',
      entry_id: this.$route.params.entry_id
    }
  },
  methods: {
    run: function() {
      var json = {"text":this.entry_id, "key":this.pass}
      this.$http.post('/api/entry/delete', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        res = response.body.substring(1,response.body.length-1);
        console.log(res)
        if (res=='success') {
          this.status = 'deleted'
          setTimeout(function(){ 
            router.push('/');  
          }, 3000);
        }
      })
    }
  }
})

var DeleteUser = Vue.extend({
  template: '#del-account',
  data: function() {
    return {
      pass: '',
      token: getCookie("token"),
      status: '',
      err: ''
    }
  },
  methods: {
    run: function() {
      var json = {"text":this.pass, "key":this.pass}
      this.$http.post('/api/user/delete', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        res = response.body.substring(1,response.body.length-1);
        if (res=='success') {
          this.status = 'deleted'
          document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
          setTimeout(function(){ 
            router.push('/');  
          }, 3000);
        } else {
          this.status = res
        }
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
        this.$http.get('/api/entry/list', {headers: {'Content-Type':'application/json', 'Authorization': 'Bearer '+getCookie("token")}}).then(response => {

        res = response.body;
        console.log(res)
        if (res == "this token is not authorized for this content") {
          
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
      formatDateTime: function(datetime) {
          d = new Date(datetime)
  
  var datestring = ("0" + d.getDate()).slice(-2) + "-" + ("0"+(d.getMonth()+1)).slice(-2) + "-" +
    d.getFullYear() + " " + ("0" + d.getHours()).slice(-2) + ":" + ("0" + d.getMinutes()).slice(-2);

  return datestring;
      },
      logout: function(){
        document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
        router.push('/login');  
      }
  },
  computed : {
    data2: function () {
      var self = this;
      return self.data.filter(function (product) {
        return self.formatDateTime(product.Day).indexOf(self.searchKey) !== -1
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
        oldUsername: '',
        has_error: false,

    }
  },
  methods: {
    createUser: function() {
      var user = this.user;
      var json = {"username": user.username, "password":user.password}
      var res = "";
      this.$http.post('/api/user/register', json, {
        headers: {
        'Content-Type': 'application/json'
        }
      }).then(response => {
        res = response.body;
        this.oldUsername = json.username;
        if (res=="\"success\""){
          this.err = "success"
          this.isDisabled = true
          setTimeout(function(){ 
            router.push('/login');  
          }, 3000);
          
        } else if (res=="\"already registered username\""){
          this.has_error = true;
          this.err = "alreadyuser"

        }
      })
      
      
    }
  }
});


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
        err2: '',
        status: ''
        // entryId: this.$route.query.entry_id
    }
  },
  components: {
        editor: summernoteComponent
    },
  created: function () {
    this.fetchData();
  },
  methods: {
    fetchData: function() {
      var json = {"text":this.$route.params.entry_id, "key":""}
      
      this.$http.post('/api/entry/view', json, {
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
      this.$http.post('/api/entry/decrypt', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        console.log(response)
        this.text = response.body;
        this.decrypted = true
        var markupStr = $('#summerasdnote').summernote('code');
        console.log(markupStr)


      })
    },

    update: function(){
      key = this.$refs[this.content.Id].value;
      idof = this.content.Id
      var json = {"text":this.text, "key":key, "idof":idof}
      this.$http.post('/api/entry/edit', json, {
        headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+getCookie("token")
        }
      }).then(response => {
        this.status = response.body.substring(1,response.body.length-1)

        if (this.status == 'success') {
          setTimeout(function(){ 
            router.push({ name: 'view', params: { entry_id: idof }});  
          }, 3000);
        }
      })

    }
  }
})


var Login = Vue.extend({
  template: '#login',
  data: function () {
    return {
        user: {username: '', password: ''},
        err: '',
        isDisabled: false,
        token: getCookie("token"),//getFromLocalStorage("Token")
        userhas_error: {'has-error': false},
        passhas_error: {'has-error': false},

    }
  },
  methods: {
    logInUser: function() {
      var user = this.user;
      var json = {"username": user.username, "password":user.password}
      var res = "";
      this.$http.post('/api/user/login', json, {
        headers: {
        'Content-Type': 'application/json'
        }
      }).then(response => {
        res = response.body;
        if (res=="user doesn't exist"){
          this.err = "doesn't exist";
          this.userhas_error['has-error'] = true
        } else if (res=="password is not true"){
          this.err = "incorrectpass";
          this.passhas_error['has-error'] = true
        } else if (res=="Error while".substring(0, 10) || res.substring(0, 10) == "everything is ".substring(0, 10)){
          this.err = "othererr";
          this.userhas_error['has-error'] = true
          this.passhas_error['has-error'] = true

        } else {
          this.err = "success";
          this.userhas_error['has-error'] = false
          this.passhas_error['has-error'] = false
          
          document.cookie = "token=" + res;

          // TODO:
          // need to implement expiration date for cookie
          // document.cookie = "token=" + res.substring(1, res.length-1) + "; expires=" + date +  "";


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
        text: '',
        key: '',
        token: getCookie("token"),
    }
  },
  components: {
        editor: summernoteComponent
    },
  methods: {
    postEntry: function() {
      var post = this.post;
      var json = {"text": this.text, "key":this.key}
      var res = "";
      this.$http.post('/api/entry/post', json, {
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
    {path: '/change-password', component: changePassword, name: 'change-password'},
    {path: '/view/:entry_id', component: View, name: 'view'},
    // {path: '/edit/:entry_id', component: Edit, name: 'edit'},
    {path: '/delete/:entry_id', component: DeleteEntry, name: 'del-entry'},
    {path: '/delete', component: DeleteUser, name: 'del-account'},
  ]
});

new Vue({
  el: '#app',
  router: router,
  template: '<router-view></router-view>'
});

