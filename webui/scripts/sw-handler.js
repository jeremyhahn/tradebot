
var REDDT_URL_CACHE = 'reddit-url-cache-v1';

// https://developers.google.com/web/fundamentals/instant-and-offline/offline-cookbook/
self.addEventListener('fetch', function(event) {
	event.respondWith(
		caches.open(REDDT_URL_CACHE).then(function(cache) {
			return fetch(event.request).then(function(response) {
				cache.put(event.request, response.clone());
				return response;
			});
		})
	);
});



