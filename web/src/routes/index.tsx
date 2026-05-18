import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
	component: Home,
});

function Home() {
	return (
		<div className="p-4">
			<h1 className="text-lg font-semibold">ZeroBudget</h1>
		</div>
	);
}
