import { AppLayout } from "@heroui-pro/react";
import { TanStackDevtools } from "@tanstack/react-devtools";
import type { QueryClient } from "@tanstack/react-query";
import {
	createRootRouteWithContext,
	HeadContent,
	Outlet,
	Scripts,
	useRouter,
} from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import TanStackQueryDevtools from "../infra/tanstack-query/devtools";
import { AppSidebar } from "../components/app-sidebar";
import appCss from "../styles.css?url";

interface MyRouterContext {
	queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
	head: () => ({
		meta: [
			{
				charSet: "utf-8",
			},
			{
				name: "viewport",
				content: "width=device-width, initial-scale=1",
			},
			{
				title: "ZeroBudget",
			},
		],
		links: [
			{
				rel: "stylesheet",
				href: appCss,
			},
		],
	}),
	shellComponent: RootDocument,
	component: RootLayout,
});

function RootLayout() {
	const router = useRouter();

	return (
		<AppLayout
			sidebarCollapsible="icon"
			navigate={(href) => router.navigate({ to: href })}
			sidebar={<AppSidebar />}
		>
			<Outlet />
		</AppLayout>
	);
}

function RootDocument({ children }: { children: React.ReactNode }) {
	return (
		<html
			lang="en"
			className="light"
			data-theme="light"
			suppressHydrationWarning
		>
			<head>
				<HeadContent />
			</head>
			<body suppressHydrationWarning>
				{children}
				<TanStackDevtools
					config={{
						position: "bottom-right",
					}}
					plugins={[
						{
							name: "Tanstack Router",
							render: <TanStackRouterDevtoolsPanel />,
						},
						TanStackQueryDevtools,
					]}
				/>
				<Scripts />
			</body>
		</html>
	);
}
