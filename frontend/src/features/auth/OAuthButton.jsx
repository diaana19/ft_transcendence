function OAuthButton() {
    const handleGithubLogin = () => {
        // Relative path so the browser hits the same origin (the edge proxy on
        // https://localhost:3000 in dev). The backend will redirect to GitHub
        // and back to /api/auth/oauth/github/callback, which sets the auth_token
        // cookie and redirects to "/".
        window.location.href = '/api/auth/oauth/github/login'
    }

    return (
        <button
            type="button"
            onClick={handleGithubLogin}
            className="w-full border border-gray-300 hover:bg-gray-100 text-black font-semibold py-3 rounded-full transition-colors flex items-center justify-center gap-3"
        >
            <img
                src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"
                alt="GitHub"
                className="w-5 h-5"
            />
            Continue with GitHub
        </button>
    )
}

export default OAuthButton
