# Checkpost

[![go-build](https://github.com/humanbeeng/checkpost/actions/workflows/go.yml/badge.svg)](https://github.com/humanbeeng/checkpost/actions/workflows/go.yml) [![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

[Checkpost](https://checkpost.io) is an open-source, real-time webhook inspection platform that allows you to easily inspect and debug incoming HTTP requests. With Checkpost, you get your own branded subdomain for free, eliminating the need to remember random UUIDs.

🚧 Project is still under construction ([track here](https://github.com/humanbeeng/checkpost/issues/100)). v1 will be released by first half of July.

### Planned Features
- [ ]  Mock server.
- [ ]  Local server via port forwarding (similar to ngrok).
- [ ]  File attachments support.
- [ ]  See through proxy middlewares with edit outgoing requests support.
- [ ]  Team plan with collaboration support across multiple URLs.
- [ ]  URL stats.

### Getting Started
1. Visit [Checkpost.io](https://checkpost.io) and sign up for a free account.
2. Create a new custom subdomain for your project.
3. Use your custom subdomain URL to start inspecting incoming HTTP requests.

### License
Checkpost.io is released under the [AGPL v3 License](https://www.gnu.org/licenses/agpl-3.0).

### Tech 
- Frontend app is built using Sveltekit. 
- Server is written in Go + Fiber.
- Database + File storage: Postgres (Supabase)
- Hosting: Cloudflare Page, Railway.app

### Contributing
Contributions to Checkpost.io are welcome! If you encounter any issues or have suggestions for improvements, please open an issue on the [repository](https://github.com/humanbeeng/checkpost/issues). If you'd like to contribute code, follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive commit messages.
4. Push your changes to your forked repository.
5. Submit a pull request to the main repository.

Please ensure that your code follows the project's coding standards and must includes appropriate tests.

#### Support
If you have any questions, suggestions or feedback, just reach out to me directly [nithin@checkpost.io](mailto:nithin@checkpost.io)  or open an issue with appropriate label.

[Terms](https://checkpost.io/terms) [Privacy](https://checkpost.io/privacy)
