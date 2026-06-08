/*
 ** File: Button.jsx
 ** Description: Reusable button component with multiple variants
 ** Responsibilities:
 ** - Render button with consistent styling
 ** - Support primary, secondary, danger and ghost variants
 ** - Support loading and disabled states
 */

import { colors } from '../../styles/tokens'

function Button({
    children,
    onClick,
    variant = 'primary',
    size = 'md',
    disabled = false,
    loading = false,
    fullWidth = false,
    type = 'button',
}) {
    const base =
        'inline-flex items-center justify-center font-bold rounded-full transition-colors disabled:opacity-50'

    const variants = {
        primary: 'bg-blue-400 hover:bg-blue-500 text-white',
        secondary: 'border border-gray-300 hover:bg-gray-100 text-gray-900',
        danger: 'border border-red-300 text-red-500 hover:bg-red-50',
        ghost: 'text-blue-400 hover:bg-blue-50',
        black: 'bg-black hover:bg-gray-800 text-white',
    }

    const sizes = {
        sm: 'px-3 py-1 text-xs',
        md: 'px-4 py-1.5 text-sm',
        lg: 'px-6 py-2 text-base',
    }

    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled || loading}
            className={`${base} ${variants[variant]} ${sizes[size]} ${fullWidth ? 'w-full' : ''}`}
        >
            {loading ? '...' : children}
        </button>
    )
}

export default Button
