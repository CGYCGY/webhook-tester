package handler

import "html/template"

var loginTmpl = template.Must(template.New("login").Parse(`<!DOCTYPE html>
<html lang="en" class="h-full">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - Webhook Tester</title>
    <link rel="stylesheet" href="/static/css/output.css">
</head>
<body class="h-full bg-gray-50 dark:bg-gray-900">
    <div class="flex min-h-full items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
        <div class="w-full max-w-md space-y-8">
            <div>
                <h1 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
                    Webhook Tester
                </h1>
                <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
                    Sign in to your account
                </p>
            </div>

            {{if .Error}}
            <div class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
                <p class="text-sm text-red-700 dark:text-red-400">{{.Error}}</p>
            </div>
            {{end}}

            <form class="mt-8 space-y-6" method="POST" action="/login">
                <div class="space-y-4 rounded-md shadow-sm">
                    <div>
                        <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                            Email address
                        </label>
                        <input
                            id="email"
                            name="email"
                            type="email"
                            autocomplete="email"
                            required
                            class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
                            placeholder="you@example.com"
                        >
                    </div>
                    <div>
                        <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                            Password
                        </label>
                        <input
                            id="password"
                            name="password"
                            type="password"
                            autocomplete="current-password"
                            required
                            class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
                            placeholder="••••••••"
                        >
                    </div>
                </div>

                <div>
                    <button
                        type="submit"
                        class="group relative flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 transition-colors"
                    >
                        Sign in
                    </button>
                </div>
            </form>
        </div>
    </div>
</body>
</html>
`))

type loginData struct {
	Error string
}
