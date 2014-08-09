var Backbone = require('backbone');
var $ = require('jquery');
var _ = require('underscore');
Backbone.$ = $;

var PomModel = Backbone.Model.extend({
	urlRoot: '/pom',
	idAttribute: 'uuid',
	defaults: function() {
		return {seconds: 0, state: 'stopped'};
		console.log(this);
	},
	sync: function(method, model, options) {
		if (options.action !== undefined) {
			options.url = [model.urlRoot,options.action,model.id].join("/")
			method = "create"
		}
		return Backbone.sync(method, model, options);
	},
	start: function() {
		var self = this;
		this.set('state', 'started');
		this.counter = setInterval(function() {
			var update = self.get('seconds')+1;
			self.set('seconds', update);
		}, 1000);
		// Post to api, set url
		this.save({}, {wait: true});
	},
	stop: function() {
		this.set('state', 'stopped');
		this.set('seconds', 0);
		this.save({}, {action: "stop"})
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
	initialize: function(opts) {
		this.gun = opts.gun;
		this.timer = this.$('[role=timer]');
		this.model.on('change:seconds', this.render, this)
		this.model.on('sync', function(){
			this.gun.trigger('pom:start', this.model);
		}, this);
		this.counter = this.$('[role=counter]');
		this.render();
	},
	render: function() {
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

var app = {};
app.gun = _.extend({}, Backbone.Events);

console.log("work fucker");
var Workspace = Backbone.Router.extend({
	routes: {
		'': 'dashboard',
		'p/:uuid': 'pom'
	},
	initialize: function() {
		console.log("router start");
	},
	dashboard: function() {
		console.log("on dashboard");
		app.dashboard = new BlinkPomo({model: new PomModel(), gun: app.gun});
	},
	pom: function(uuid) {
		console.log("uuid", uuid);
	}
});

app.router = new Workspace();

app.gun.on('pom:start', function(pom) {
	app.router.navigate("/p/" + pom.id);
});

console.log(
	"App?",
	Backbone.history.start({pushState: true, silent: false})
);
