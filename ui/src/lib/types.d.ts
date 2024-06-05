export type User = {
	id: number;
	avatar_url: string;
	name: string;
	email: string;
	plan: string;
};

export type Plan = 'free' | 'guest' | 'hobby' | 'pro';

export type Endpoint = {
	endpoint: string;
	plan: Plan;
	expires_at: string;
};

export type HttpMethod =
	| 'get'
	| 'post'
	| 'delete'
	| 'patch'
	| 'put'
	| 'connect'
	| 'head'
	| 'options';

export type Request = {
	endpoint: string;
	path: string;
	content: string;
	method: HttpMethod;
	uuid: string;
	source_ip: string;
	content_size: number;
	response_code: number;
	headers: Object;
	query_params: Object;
	created_at: string;
	expires_at: string;
};

export type UrlHistory = {
	requests: Request[];
};
