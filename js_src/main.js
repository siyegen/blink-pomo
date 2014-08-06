var Backbone = require('backbone');
var $ = require('jquery');
Backbone.$ = $;

var PomModel = Backbone.Model.extend({
	defaults: function() {
		return {seconds: 0, state: 'stopped'};
		console.log(this);
	},
	start: function() {
		var self = this;
		this.set('state', 'started');
		this.counter = setInterval(function() {
			var update = self.get('seconds')+1;
			self.set('seconds', update);
		}, 1000);
	},
	stop : function() {
		this.set('state', 'stopped');
		clearInterval(this.counter);
	},
	formatCounter: function() {
		var min = ('0'+Math.floor(this.get('seconds') / 60)).slice(-2);
		var secs = ('0'+this.get('seconds') % 60).slice(-2);
		return min + ':' + secs;
	}
});

var BlinkPomo = Backbone.View.extend({
	el: $('[role=app]'),
	events: {
		'click [role=start-pom]': 'start',
		'click [role=stop-pom]': 'stop',
	},
	initialize: function() {
		// _.bindAll(this, 'render');
		this.timer = this.$('[role=timer]');
		this.model.on('change:seconds', this.render, this)
		this.counter = this.$('[role=counter]');
		this.render();
	},
	render: function() {
		console.log('called render');
		this.counter.text(this.model.formatCounter());
	},
	start: function(){
		console.log('starting pom');
		this.model.start();
		this.render();
	},
	stop: function(){
		console.log('stopping pom');
		this.model.stop();
		this.render();
	}
});

var app = new BlinkPomo({model: new PomModel()});
