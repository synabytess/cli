// NAME: Christian Spotify
// AUTHOR: khanhas
// DESCRIPTION: Auto skip explicit songs. Toggle in Profile menu.

/// <reference path="../globals.d.ts" />

(async function ChristianSpotify() {
	if (!skidify.LocalStorage) {
		setTimeout(ChristianSpotify, 1000);
		return;
	}
	await new Promise((res) => skidify.Events.webpackLoaded.on(res));

	let isEnabled = skidify.LocalStorage.get("ChristianMode") === "1";

	new skidify.Menu.Item("Christian mode", isEnabled, (self) => {
		isEnabled = !isEnabled;
		skidify.LocalStorage.set("ChristianMode", isEnabled ? "1" : "0");
		self.setState(isEnabled);
	}).register();

	skidify.Player.addEventListener("songchange", () => {
		if (!isEnabled) return;
		const data = skidify.Player.data || skidify.Queue;
		if (!data) return;

		const isExplicit = data.item.metadata.is_explicit;
		if (isExplicit === "true") {
			skidify.Player.next();
		}
	});
})();
