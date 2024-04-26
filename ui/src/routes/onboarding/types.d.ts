export type User = {
	id: number;
	name: string;
	email: string;
	plan: Plan;
};

export type Plan = 'free' | 'guest' | 'hobby' | 'pro';

export type GenerateUrlRequest = {
	endpoint: string;
};

export type UserEndpointsResponse = {
	endpoints: Endpoint[];
};

export type Endpoint = {
	endpoint: string;
	plan: Plan;
	expiresAt: string;
};
