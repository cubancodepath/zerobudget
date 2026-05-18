import ky from "ky";

const configuredApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, "");
const apiBaseUrl = configuredApiBaseUrl ?? "http://localhost:8080";

export const http = ky.create({
	baseUrl: `${apiBaseUrl}/`,
	headers: {
		"Content-Type": "application/json",
	},
	timeout: 10000,
});
