import { Button, Kbd } from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";
import { Chart } from "../components/AreaChart";

export const Route = createFileRoute("/")({ component: Home });

function Home() {
	return (
		<div className="p-4">
			<Button variant="primary">Click me</Button>
			<Chart />
			<div className="flex items-center gap-4">
				<Kbd>
					<Kbd.Abbr keyValue="command" />
					<Kbd.Content>K</Kbd.Content>
				</Kbd>
				<Kbd>
					<Kbd.Abbr keyValue="shift" />
					<Kbd.Content>P</Kbd.Content>
				</Kbd>
				<Kbd>
					<Kbd.Abbr keyValue="ctrl" />
					<Kbd.Content>C</Kbd.Content>
				</Kbd>
				<Kbd>
					<Kbd.Abbr keyValue="option" />
					<Kbd.Content>D</Kbd.Content>
				</Kbd>
			</div>
		</div>
	);
}
