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
		if (this.isNew()) {
			console.log("Can't start without a chrono");
			return
		}
		var self = this;
		this.save({'state': 'starting'}, {wait: true});
		this.counter = setInterval(function() {
			var update = self.get('seconds')+1;
			self.set('seconds', update);
		}, 1000);
	},
	stop: function() {
		if (this.isNew()) {
			console.log("Can't stop without a chrono");
			return
		}
		this.save({'state': 'stopping', 'seconds': 0});
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
		'click [role=create-pom]': 'create',
	},
	initialize: function(opts) {
		console.log("id of model", this.model.id);
		this.gun = opts.gun;
		this.timer = this.$('[role=timer]');
		this.model.on('change:seconds', this.render, this)
		this.model.on('sync', function(){
			this.gun.trigger('pom:start', this.model);
		}, this);
		this.counter = this.$('[role=counter]');
		this.ui_create = this.$('[role=create-pom]');
		this.render();
	},
	render: function() {
		this.counter.text(this.model.formatCounter());
	},
	create: function() {
		console.log('creating pom');
		this.ui_create.hide();
		this.model.save({});
	},
	start: function() {
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

// Used to create a pom
// var Chrono = Backbone.View.extend({

var app = {};
app.gun = _.extend({}, Backbone.Events);

var Workspace = Backbone.Router.extend({
	routes: {
		'chrono/': 'chrono',
		'chrono/:uuid': 'pom'
	},
	initialize: function() {
		console.log("router start");
	},
	chrono: function() {
		console.log("on dashboard");
		app.chrono = new BlinkPomo({model: new PomModel(), gun: app.gun});
	},
	pom: function(uuid) {
		console.log("uuid", uuid);
		app.chrono = new BlinkPomo({model: new PomModel({id:uuid}), gun: app.gun});
		window.chrono = app.chrono;
	}
});

app.router = new Workspace();

app.gun.on('pom:start', function(pom) {
	app.router.navigate("/chrono/" + pom.id);
});

console.log(
	"App?",
	Backbone.history.start({pushState: true, silent: false})
);
