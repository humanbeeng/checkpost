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
