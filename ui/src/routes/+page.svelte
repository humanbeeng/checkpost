<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
</script>

<header
	class="antialiased sticky top-0 border-b border-black/10 w-full bg-gray-100/10 backdrop-blur-3xl"
>
	<div class="container flex h-12 max-w-screen-2xl items-center justify-between">
		<a href="/" class="font-medium">Checkpost</a>
		<nav class="flex items-center gap-6 text-sm">
			<a href="/pricing" class="hover:underline">Pricing</a>
			<a href="/about" class="hidden md:flex hover:underline">About</a>
			<a href="/admin/dashboard" class="hover:underline">Dashboard</a>
			<Button class="h-8 hidden border-black md:flex" href="/auth/github" variant="default"
				>Sign in</Button
			>
		</nav>
	</div>
</header>
<div class=" w-full min-h-full antialiased">
	<svg
		xmlns="http://www.w3.org/2000/svg"
		version="1.1"
		xmlns:xlink="http://www.w3.org/1999/xlink"
		preserveAspectRatio="none"
		class="absolute inset-0 -z-10 h-full w-full stroke-gray-200 [mask-image:radial-gradient(100%_100%_at_top_right,white,transparent)]"
	>
		<g mask="url(&quot;#SvgjsMask1028&quot;)" fill="none">
			<use xlink:href="#SvgjsSymbol1035" x="0" y="0"></use>
			<use xlink:href="#SvgjsSymbol1035" x="0" y="720"></use>
			<use xlink:href="#SvgjsSymbol1035" x="720" y="0"></use>
			<use xlink:href="#SvgjsSymbol1035" x="720" y="720"></use>
			<use xlink:href="#SvgjsSymbol1035" x="1440" y="0"></use>
			<use xlink:href="#SvgjsSymbol1035" x="1440" y="720"></use>
		</g>
		<defs>
			<mask id="SvgjsMask1028">
				<rect width="1920" height="1080" fill="#ffffff"></rect>
			</mask>
			<path d="M-1 0 a1 1 0 1 0 2 0 a1 1 0 1 0 -2 0z" id="SvgjsPath1031"></path>
			<path d="M-3 0 a3 3 0 1 0 6 0 a3 3 0 1 0 -6 0z" id="SvgjsPath1034"></path>
			<path d="M-5 0 a5 5 0 1 0 10 0 a5 5 0 1 0 -10 0z" id="SvgjsPath1030"></path>
			<path d="M2 -2 L-2 2z" id="SvgjsPath1033"></path>
			<path d="M6 -6 L-6 6z" id="SvgjsPath1029"></path>
			<path d="M30 -30 L-30 30z" id="SvgjsPath1032"></path>
		</defs>
		<symbol id="SvgjsSymbol1035">
			<use xlink:href="#SvgjsPath1029" x="30" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="30" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="30" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="30"
				y="210"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="30"
				y="270"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1033" x="30" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="30" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="30" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="30" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="30" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="30" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="30" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="90"
				y="30"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1030" x="90" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="90" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="90" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="90" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="90" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="90" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="90" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="90" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="90" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="90" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="90" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="150" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="150" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="150" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="150" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="150" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="150" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="150" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="150" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="150" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="150"
				y="570"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1029" x="150" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="150" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="210" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="210" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="210" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="210" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="210" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="210" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="210" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="210" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="210" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="210" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="210" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="210" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="270" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="270" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="270"
				y="150"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1030" x="270" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="270" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="270" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="270" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="270" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="270" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="270" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="270" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="270" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="330" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="330" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="330" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="330" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="330" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="330" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="330" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="330" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="330" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="330" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="330" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="330" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="390" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="390" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="390" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="390" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="390" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="390" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="390" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="390" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="390" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="390" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="390" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="390" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="450" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="450" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="450" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="450" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="450" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="450" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="450" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="450" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="450" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="450" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="450" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="450" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="510"
				y="30"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1029" x="510" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="510" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="510" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="510" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="510" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="510" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="510" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="510" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="510"
				y="570"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1031" x="510" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="510" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="570" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="570" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="570" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="570" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="570" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="570" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="570" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="570" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="570" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="570" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="570" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="570" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="630" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="630" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="630" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="630" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="630" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1031" x="630" y="330" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="630" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="630" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="630" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="630" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="630" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="630" y="690" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="690" y="30" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="690" y="90" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1033" x="690" y="150" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1030" x="690" y="210" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="690" y="270" stroke="rgba(217, 217, 217, 1)"></use>
			<use
				xlink:href="#SvgjsPath1032"
				x="690"
				y="330"
				stroke="rgba(217, 217, 217, 1)"
				stroke-width="3"
			></use>
			<use xlink:href="#SvgjsPath1031" x="690" y="390" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="690" y="450" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="690" y="510" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="690" y="570" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1034" x="690" y="630" stroke="rgba(217, 217, 217, 1)"></use>
			<use xlink:href="#SvgjsPath1029" x="690" y="690" stroke="rgba(217, 217, 217, 1)"></use>
		</symbol>
	</svg>

	<div class=" grid grid-cols-1 mx-6 md:min-h-screen md:grid-cols-2 md:max-w-7xl lg:mx-auto">
		<div class="w-auto flex my-32 flex-col md:justify-center">
			<h1
				class="text-4xl antialiased font-bold leading-tight select-none tracking-tighter bg-gradient-to-r text-left lg:text-6xl lg:text-left lg:leading-[1.1]"
			>
				Inspect incoming HTTP requests in realtime
			</h1>

			<section>
				<h3
					class="font-light text-xl opacity-60 pt-2 select-none leading-10 md:py-4 md:leading-10 text-left"
				>
					Within seconds, you'll get to
					<ul class="list-disc mx-5 font-normal text-lg">
						<li>access a free endpoint</li>
						<li>view incoming requests in realtime</li>
					</ul>
				</h3>
			</section>

			<Button class="w-full h-11 bg-foreground my-4 lg:w-1/2">Coming soon.</Button>
		</div>
		<!--Dashboard image-->
		<div class="hidden md:w-auto md:block"></div>
	</div>

	<section class=" h-screen border-t">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="absolute w-full h-full -z-10 [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]"
			><rect fill="#d9d9d9" /><defs
				><rect stroke="#ffffff" stroke-width="0.5" width="1" height="1" id="s" /><pattern
					id="a"
					width="3"
					height="3"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><use fill="#fcfcfc" href="#s" y="2" /><use fill="#fcfcfc" href="#s" x="1" y="2" /><use
						fill="#fafafa"
						href="#s"
						x="2"
						y="2"
					/><use fill="#fafafa" href="#s" /><use fill="#f7f7f7" href="#s" x="2" /><use
						fill="#f7f7f7"
						href="#s"
						x="1"
						y="1"
					/></pattern
				><pattern
					id="b"
					width="7"
					height="11"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#f5f5f5"
						><use href="#s" /><use href="#s" y="5" /><use href="#s" x="1" y="10" /><use
							href="#s"
							x="2"
							y="1"
						/><use href="#s" x="2" y="4" /><use href="#s" x="3" y="8" /><use
							href="#s"
							x="4"
							y="3"
						/><use href="#s" x="4" y="7" /><use href="#s" x="5" y="2" /><use
							href="#s"
							x="5"
							y="6"
						/><use href="#s" x="6" y="9" /></g
					></pattern
				><pattern
					id="h"
					width="5"
					height="13"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#f5f5f5"
						><use href="#s" y="5" /><use href="#s" y="8" /><use href="#s" x="1" y="1" /><use
							href="#s"
							x="1"
							y="9"
						/><use href="#s" x="1" y="12" /><use href="#s" x="2" /><use href="#s" x="2" y="4" /><use
							href="#s"
							x="3"
							y="2"
						/><use href="#s" x="3" y="6" /><use href="#s" x="3" y="11" /><use
							href="#s"
							x="4"
							y="3"
						/><use href="#s" x="4" y="7" /><use href="#s" x="4" y="10" /></g
					></pattern
				><pattern
					id="c"
					width="17"
					height="13"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#f2f2f2"
						><use href="#s" y="11" /><use href="#s" x="2" y="9" /><use href="#s" x="5" y="12" /><use
							href="#s"
							x="9"
							y="4"
						/><use href="#s" x="12" y="1" /><use href="#s" x="16" y="6" /></g
					></pattern
				><pattern
					id="d"
					width="19"
					height="17"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#ffffff"
						><use href="#s" y="9" /><use href="#s" x="16" y="5" /><use href="#s" x="14" y="2" /><use
							href="#s"
							x="11"
							y="11"
						/><use href="#s" x="6" y="14" /></g
					><g fill="#efefef"
						><use href="#s" x="3" y="13" /><use href="#s" x="9" y="7" /><use
							href="#s"
							x="13"
							y="10"
						/><use href="#s" x="15" y="4" /><use href="#s" x="18" y="1" /></g
					></pattern
				><pattern
					id="e"
					width="47"
					height="53"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#FFFFFF"
						><use href="#s" x="2" y="5" /><use href="#s" x="16" y="38" /><use
							href="#s"
							x="46"
							y="42"
						/><use href="#s" x="29" y="20" /></g
					></pattern
				><pattern
					id="f"
					width="59"
					height="71"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#FFFFFF"
						><use href="#s" x="33" y="13" /><use href="#s" x="27" y="54" /><use
							href="#s"
							x="55"
							y="55"
						/></g
					></pattern
				><pattern
					id="g"
					width="139"
					height="97"
					patternUnits="userSpaceOnUse"
					patternTransform="scale(7.25) translate(-862.07 -646.55)"
					><g fill="#FFFFFF"
						><use href="#s" x="11" y="8" /><use href="#s" x="51" y="13" /><use
							href="#s"
							x="17"
							y="73"
						/><use href="#s" x="99" y="57" /></g
					></pattern
				></defs
			><rect fill="url(#a)" width="100%" height="100%" /><rect
				fill="url(#b)"
				width="100%"
				height="100%"
			/><rect fill="url(#h)" width="100%" height="100%" /><rect
				fill="url(#c)"
				width="100%"
				height="100%"
			/><rect fill="url(#d)" width="100%" height="100%" /><rect
				fill="url(#e)"
				width="100%"
				height="100%"
			/><rect fill="url(#f)" width="100%" height="100%" /><rect
				fill="url(#g)"
				width="100%"
				height="100%"
			/></svg
		>
		<h1 class="text-3xl font-bold tracking-tight text-center py-10 md:py-12">Plans</h1>
		<div class="w-full px-6 lg:px-32 grid grid-cols-1 gap-6 md:grid-cols-3">
			<!-- Guest -->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Guest</Card.Title>
					<Card.Description>Non signed in users</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>One random endpoint. <br>Example: <a href="/"><u>https://xyzxa.checkpost.io</u></a> </li>
						<li>Expires in 1 day.</li>
						<li>Request history retained for 1 day.</li>
						<li>Request payload limit 32Kb.</li>
						<li>Rate limit 2 rps.</li>
					</ul>
				</Card.Content>
			</Card.Root>

			<!-- Free -->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Free</Card.Title>
					<Card.Description>Signed in users</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>One random endpoint. <br>Example: <a href="/"><u>https://xyzxa.checkpost.io</u></a> </li>
						<li>Expires in 2 days.</li>
						<li>Request history retained for 2 days.</li>
						<li>Request payload limit 32Kb.</li>
						<li>Rate limit: 3 rps.</li>
						<li>Discord Access - Community channel.</li>
					</ul>
				</Card.Content>
				<!-- <Card.Footer> -->
				<!-- 	<Button class="w-full">Generate endpoint</Button> -->
				<!-- </Card.Footer> -->
			</Card.Root>

			<!--Pro-->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="underline text-xl">Pro</Card.Title>
					<Card.Description>$5/month | $55/year</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>10 custom endpoints. <br><a href="/"><u>Example: https://yourname.checkpost.io</u></a></li>
						<li>Request history is stored upto 180 days.</li>
						<li>10 Web CRONs.</li>
						<li>Request payload limit 510Kb.</li>
						<li>Password protected endpoints.</li>
						<li>File attachments support.</li>
						<li>Port forwarding for local servers.</li>
						<li>1 proxy middleware.</li>
						<li>Rate limit: 30 rps.</li>
						<li>Discord Access - Private Channel</li>
					</ul>
				</Card.Content>
				<Card.Footer>
					<Button class="w-full bg-black underline">Upgrade - Coming soon</Button>
				</Card.Footer>
			</Card.Root>

			<!--Enterprise-->
			<!-- <Card.Root> -->
			<!-- 	<Card.Header> -->
			<!-- 		<Card.Title tag="h1" class="text-xl">Enterprise (Coming soon)</Card.Title> -->
			<!-- 		<Card.Description>Enterprise (Coming soon)</Card.Description> -->
			<!-- 	</Card.Header> -->
			<!-- <Card.Content> -->
			<!-- 	<ul class="list-disc mx-5"> -->
			<!-- 		<li>Inspect incoming HTTP requests.</li> -->
			<!-- 		<li>1 endpoint</li> -->
			<!-- 		<li>Request history retained for 2 days</li> -->
			<!-- 		<li>Expires in 4h.</li> -->
			<!-- 		<li>Request payload limit 32Kb.</li> -->
			<!-- 		<li>Rate limit: 10 rps.</li> -->
			<!-- 		<li>Discord Access</li> -->
			<!-- 	</ul> -->
			<!-- </Card.Content> -->
			<!-- <Card.Footer> -->
			<!-- 	<Button class="w-full">Generate endpoint</Button> -->
			<!-- </Card.Footer> -->
			<!-- </Card.Root> -->
		</div>
		<footer class="py-4 md:px-8 md:py-0 bottom-0">
			<div class="container flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
				<div class="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
					<p class="text-center text-sm leading-loose text-muted-foreground md:text-left">
						Built by <a
							href="https://twitter.com/nithinrajx"
							target="_blank"
							rel="noreferrer"
							class="font-medium underline underline-offset-4"
							data-svelte-h="svelte-18aekvv">Nithin Raj</a
						>.
					</p>
				</div>
			</div>
		</footer>
	</section>
</div>
