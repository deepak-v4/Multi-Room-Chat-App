new Vue({
  el: '#app',

  data: {
      ws: null, // Our websocket
      newMsg: '', // Holds new messages to be sent to the server
      chatContent: '', // A running list of chat messages displayed on the screen
      //email: null, // Email address used for grabbing an avatar
      username: null, // Our username
      roomid:'Room101',// default room id
      joined: false // True if email and username have been filled in
  },

  created: function() {
      var self = this;
      this.ws = new WebSocket('ws://' + window.location.host + '/ws');
      this.ws.addEventListener('message', function(e) {
          var msg = JSON.parse(e.data);
          self.chatContent += '<div class="chip">'
                  + '<img src="' + self.gravatarURL(msg.username) + '">' // Avatar
                  + msg.username
              + '</div>'
              + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

          var element = document.getElementById('chat-messages');
          element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
      });
  },

  methods: {
      send: function () {

        console.log(this.newMsg)
          if (this.newMsg != '') {
              this.ws.send(
                  JSON.stringify({
                      //email: this.email,
                      username: this.username,
                      roomid:this.roomid,
                      message: $('<p>').html(this.newMsg).text() // Strip out html
                  }
              ));
              this.newMsg = ''; // Reset newMsg
          }
      },

      join: function () {
          /*if (!this.email) {
              Materialize.toast('You must enter an email', 2000);
              return
          }*/
          if (!this.username) {
              Materialize.toast('You must choose a username', 2000);
              return
          }
          //this.email = $('<p>').html(this.email).text();
          this.username = $('<p>').html(this.username).text();
          this.roomid=$('<p>').html(this.roomid).text();
          this.joined = true;
      },

      gravatarURL: function(usr) {
          return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(usr);
      }
  }
});