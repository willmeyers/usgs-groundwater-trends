import { writable } from "svelte/store";
import type { Site, Mouse, SiteCache } from "./lib/types";

export const mouse = writable<Mouse>({ x: 0, y: 0 });
export const selectedSite = writable<Site | null>(null);
export const cache = writable<SiteCache>({});
