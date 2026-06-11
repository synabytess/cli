// NAME: Auto Skip Video
// AUTHOR: khanhas
// DESCRIPTION: Auto skip video

/// <reference path="../globals.d.ts" />

(function SkipVideo() {
	skidify.Player.addEventListener("songchange", () => {
		const data = skidify.Player.data || skidify.Queue;
		if (!data) return;

		const meta = data.item.metadata;
		// Ads are also video media type so I need to exclude them out.
		if (meta["media.type"] === "video" && meta.is_advertisement !== "true") {
			skidify.Player.next();
		}
	});
})();
