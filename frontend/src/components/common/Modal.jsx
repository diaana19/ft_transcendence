function Modal({
    isOpen,
    onClose,
    title,
    children,
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    onConfirm = null,
    confirmVariant = 'primary',
    size = 'md',
}) {
    if (!isOpen) return null

    const sizes = {
        sm: 'max-w-sm',
        md: 'max-w-md',
        lg: 'max-w-lg',
    }

    const confirmStyles = {
        primary: 'bg-blue-400 hover:bg-blue-500 text-white',
        danger: 'bg-red-500 hover:bg-red-600 text-white',
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div
                className="absolute inset-0 backdrop-blur-sm"
                style={{
                    background:
                        'linear-gradient(135deg, #faf5ff80 10%, #eef2ff80 60%, #f5f7fb80 100%)',
                }}
                onClick={onClose}
            />

            <div
                className={`relative z-10 bg-white rounded-2xl w-full ${sizes[size]} shadow-xl overflow-hidden mx-4`}
            >
                <div className="flex items-center gap-4 px-4 py-3 border-b border-gray-200">
                    <button
                        onClick={onClose}
                        className="hover:bg-gray-100 rounded-full p-2 text-sm text-gray-500"
                    >
                        ✕
                    </button>
                    <h3 className="text-lg font-bold flex-1 text-gray-900">{title}</h3>
                </div>

                <div className="px-4 py-4">{children}</div>

                {onConfirm && (
                    <div className="flex justify-end gap-2 px-4 pb-4">
                        <button
                            onClick={onClose}
                            className="px-4 py-2 rounded-full border border-gray-300 text-sm font-semibold hover:bg-gray-100 transition-colors"
                        >
                            {cancelText}
                        </button>
                        <button
                            onClick={onConfirm}
                            className={`px-4 py-2 rounded-full text-sm font-bold transition-colors ${confirmStyles[confirmVariant]}`}
                        >
                            {confirmText}
                        </button>
                    </div>
                )}
            </div>
        </div>
    )
}

export default Modal
