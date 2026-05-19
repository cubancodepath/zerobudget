import { Button, Modal } from "@heroui/react";
import {
	createContext,
	type ReactNode,
	useCallback,
	useContext,
	useMemo,
	useState,
} from "react";

type ModalSize = "xs" | "sm" | "md" | "lg" | "cover" | "full";
type ModalPlacement = "auto" | "top" | "center" | "bottom";
type ModalBackdropVariant = "opaque" | "blur" | "transparent";

type ModalRenderContext = {
	close: () => void;
};

export type OpenModalOptions = {
	render?: (ctx: ModalRenderContext) => ReactNode;
	size?: ModalSize;
	placement?: ModalPlacement;
	backdrop?: ModalBackdropVariant;
	isDismissable?: boolean;
	isKeyboardDismissDisabled?: boolean;
	onOpenChange?: (isOpen: boolean) => void;
};

type ModalContextValue = {
	isOpen: boolean;
	open: (options?: OpenModalOptions) => boolean;
	close: () => void;
};

type ModalProviderProps = {
	children: ReactNode;
	defaultModalOptions?: OpenModalOptions;
};

const FALLBACK_DEFAULTS: Required<
	Pick<
		OpenModalOptions,
		| "size"
		| "placement"
		| "backdrop"
		| "isDismissable"
		| "isKeyboardDismissDisabled"
	>
> = {
	size: "md",
	placement: "center",
	backdrop: "opaque",
	isDismissable: true,
	isKeyboardDismissDisabled: false,
};

function defaultRender({ close }: ModalRenderContext) {
	return (
		<Modal.Dialog className="sm:max-w-md">
			<Modal.CloseTrigger />
			<Modal.Header>
				<Modal.Heading>Modal</Modal.Heading>
			</Modal.Header>
			<Modal.Body>
				<p className="text-sm text-muted">
					This is the default modal. Pass <code>render</code> to open a custom
					one.
				</p>
			</Modal.Body>
			<Modal.Footer>
				<Button variant="secondary" onPress={close}>
					Close
				</Button>
			</Modal.Footer>
		</Modal.Dialog>
	);
}

const ModalContext = createContext<ModalContextValue | null>(null);

export function ModalProvider({
	children,
	defaultModalOptions,
}: ModalProviderProps) {
	const [currentModal, setCurrentModal] = useState<OpenModalOptions | null>(
		null,
	);

	const close = useCallback(() => {
		setCurrentModal((prev) => {
			if (!prev) return prev;
			prev.onOpenChange?.(false);
			return null;
		});
	}, []);

	const open = useCallback(
		(options?: OpenModalOptions) => {
			let opened = false;

			setCurrentModal((prev) => {
				if (prev) return prev;

				const defaults = defaultModalOptions ?? {};
				const next: OpenModalOptions = {
					...defaults,
					...options,
					render: options?.render ?? defaults.render ?? defaultRender,
				};

				next.onOpenChange?.(true);
				opened = true;
				return next;
			});

			return opened;
		},
		[defaultModalOptions],
	);

	const value = useMemo<ModalContextValue>(
		() => ({
			isOpen: currentModal !== null,
			open,
			close,
		}),
		[currentModal, open, close],
	);

	const modalConfig = currentModal ?? defaultModalOptions;

	return (
		<ModalContext.Provider value={value}>
			{children}
			<Modal.Backdrop
				isOpen={currentModal !== null}
				onOpenChange={(isOpen) => {
					if (!isOpen) close();
				}}
				variant={modalConfig?.backdrop ?? FALLBACK_DEFAULTS.backdrop}
				isDismissable={
					modalConfig?.isDismissable ?? FALLBACK_DEFAULTS.isDismissable
				}
				isKeyboardDismissDisabled={
					modalConfig?.isKeyboardDismissDisabled ??
					FALLBACK_DEFAULTS.isKeyboardDismissDisabled
				}
			>
				<Modal.Container
					size={modalConfig?.size ?? FALLBACK_DEFAULTS.size}
					placement={modalConfig?.placement ?? FALLBACK_DEFAULTS.placement}
				>
					{currentModal?.render?.({ close })}
				</Modal.Container>
			</Modal.Backdrop>
		</ModalContext.Provider>
	);
}

export function useModal() {
	const context = useContext(ModalContext);

	if (!context) {
		throw new Error("useModal must be used within a ModalProvider");
	}

	return context;
}
