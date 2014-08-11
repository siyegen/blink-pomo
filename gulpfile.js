var gulp = require('gulp');
var browserify = require('gulp-browserify');
var concat = require('gulp-concat');

gulp.task('dev', function() {
	gulp.start('browserify');
	gulp.watch('./js_src/**/*.js', ['browserify']);
});

gulp.task('browserify', function() {
	return gulp.src(['./js_src/main.js'])
		.pipe(browserify())
        .pipe(concat('blink-pomo.js'))
		.pipe(gulp.dest('./assets/js/'));
});
