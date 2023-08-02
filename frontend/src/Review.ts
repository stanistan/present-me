export interface MaybeLinked {
  text: string;
  href?: string;
}

export interface Comment {
  number: number;
  title: MaybeLinked;
  description: string;
  code: {
    content: string;
    lang: string;
    diff?: boolean;
  };
}

export interface Link {
  label: string;
  text: string;
  href: string;
}

export interface Review {
  title: MaybeLinked;
  body: string;
  links?: Array<Link>;
  comments: Array<Comment>;
  meta: {
    params: {};
  };
}
