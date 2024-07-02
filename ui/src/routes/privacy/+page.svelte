<script>
	import Footer from '@/components/Footer.svelte';
	import Header from '@/components/Header.svelte';
</script>

<body class="bg-gray-50 w-full h-screen">
	<Header />
	<div class="flex flex-col justify-center mx-10 lg:mx-64">
		<h1 class="text-4xl font-bold tracking-tight my-16 text-center">Privacy Policy</h1>

		<div class="w-full text-md flex flex-col gap-7">
			<section>
				At <a href="https://checkpost.io" class="underline">Checkpost.io</a>, we take data privacy
				seriously. This Privacy Policy explains how we collect, use, disclose, and protect the
				personal information you provide when using our webhook and HTTP request inspection
				services.
			</section>

			<section>
				<h2 class="text-xl font-medium">Information We Collect</h2>
				<ul>
					<li>
						1. Account Information: When you create an account, we store your email address, name,
						and profile picture. This details are used only to identify users and to receive any
						communication from the service.
					</li>
					<li>
						2. Request Data: When testing HTTP requests through our service, we log the request
						details such as headers, content, IP, client data, file attachments. This data is
						associated with your account to provide the request inspection functionality.
					</li>
				</ul>
			</section>

			<section>
				<h2 class="text-xl font-medium">Use of Information</h2>
				<ul>
					<li>
						1. We use account information to provide the checkpost.io services, authenticate users,
						and send account notifications.
					</li>
					<li>
						2. Request data is used solely to display details of the HTTP requests you test through
						our service.
					</li>
					<li>
						3. We do not sell or rent personal information and request details to third parties.
					</li>
				</ul>
			</section>

			<section class="flex flex-col gap-2">
				<p>
					<span class="font-medium"> Data Retention: </span>We retain account information for as
					long as your account is active. Request data is retained in accordance with user's data
					retention plan.
				</p>

				<p>
					<span class="font-medium"> Disclosure of Information: </span>
					We may disclose personal information if we believe in good faith that such disclosure is necessary
					to comply with relevant laws or respond to subpoenas, court orders, or other legal processes.
				</p>

				<p>
					<span class="font-medium">Data Security: </span>We use industry-standard security measures
					to protect personal information from unauthorized access or disclosure. However, no method
					of Internet transmission or storage is 100% secure.
				</p>

				<p>
					<span class="font-medium"> Links to Other Sites: </span>Our website may contain links to
					third-party sites. We are not responsible for their privacy practices.
				</p>

				<p>
					<span class="font-medium">Children's Privacy: </span>Our services are not intended for
					children under 13. We do not knowingly collect personal information from children.
				</p>
				<p>
					Updates to this Policy We may update this Privacy Policy from time to time. We encourage
					you to review it periodically.
				</p>

				<p>
					Contact Us For any questions about this policy or our data practices, please contact us at
					<a class="underline text-accent-foreground" href="mailto:nithin@checkpost.io"
						>nithin@checkpost.io</a
					>
				</p>
			</section>
		</div>
	</div>

	<Footer></Footer>
</body>
