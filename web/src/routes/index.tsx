import { Button, Kbd } from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: Home });

function Home() {
	return (
		<div className="p-4">
			<Button variant="secondary">Click me</Button>
			<div className="flex items-center gap-4">
				<Kbd>
					<Kbd.Abbr keyValue="command">CMD</Kbd.Abbr>
					<Kbd.Content>K</Kbd.Content>
				</Kbd>
				<Kbd>
					<Kbd.Abbr keyValue="shift">SHIFT</Kbd.Abbr>
					<Kbd.Content>P</Kbd.Content>
				</Kbd>
				<Kbd>
					<Kbd.Abbr keyValue="ctrl">CTRL</Kbd.Abbr>
					<Kbd.Content>C</Kbd.Content>
				</Kbd>
				<Kbd>
					<Kbd.Abbr keyValue="option">OPT</Kbd.Abbr>
					<Kbd.Content>D</Kbd.Content>
				</Kbd>
			</div>
		</div>
	);
}
