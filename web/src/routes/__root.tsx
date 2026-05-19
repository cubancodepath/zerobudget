import { AppLayout } from "@heroui-pro/react";
import { useLiveQuery } from "@tanstack/react-db";
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
import { useEffect, useState } from "react";
import type {
	AccountOfflineExecutor,
	AccountsCollection,
} from "#/collections/accounts";
import { ModalProvider } from "#/components/modal";
import { AppSidebar } from "../components/app-sidebar";
import TanStackQueryDevtools from "../infra/tanstack-query/devtools";
import appCss from "../styles.css?url";

interface MyRouterContext {
	queryClient: QueryClient;
	accountsCollection: AccountsCollection;
	accountsOfflineExecutor: AccountOfflineExecutor;
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
	ssr: false,
});

function RootLayout() {
	const router = useRouter();
	const { accountsCollection, accountsOfflineExecutor } =
		Route.useRouteContext();
	const [ifOfflineReady, setIfOfflineReady] = useState(false);

	useEffect(() => {
		const isActive = true;
		void accountsOfflineExecutor.waitForInit().finally(() => {
			if (isActive) setIfOfflineReady(true);
		});
	}, [accountsOfflineExecutor]);

	const { data: accounts = [] } = useLiveQuery((q) =>
		q.from({ accounts: accountsCollection }),
	);

	return (
		<ModalProvider>
			<AppLayout
				sidebarCollapsible="icon"
				navigate={(href) => router.navigate({ to: href })}
				sidebar={<AppSidebar accounts={accounts} isReady={ifOfflineReady} />}
			>
				<Outlet />
			</AppLayout>
		</ModalProvider>
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
			<body suppressHydrationWarning suppressContentEditableWarning>
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
