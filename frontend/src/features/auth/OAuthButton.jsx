function OAuthButton() {
    const handleGithubLogin = () => {
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
