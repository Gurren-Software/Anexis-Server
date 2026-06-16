import type { ReactNode } from 'react';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

type Feature = {
  label: string;
  title: string;
  description: string;
};

const features: Feature[] = [
  {
    label: 'Storage',
    title: 'Multi-provider by design',
    description:
      'Local disks, Backblaze B2 and S3-compatible storage can live behind the same API.',
  },
  {
    label: 'Modes',
    title: 'SaaS or standalone',
    description:
      'Run multi-tenant with quotas, or deploy a private instance protected by an API key.',
  },
  {
    label: 'Ops',
    title: 'Built for self-hosting',
    description:
      'Docker-first deployment, explicit environment config and predictable production paths.',
  },
  {
    label: 'Data',
    title: 'Migration-ready',
    description:
      'Move files from cloud providers and keep backup exports close at hand.',
  },
];

const quickLinks = [
  { label: 'Install', to: '/docs/installation' },
  { label: 'Configure', to: '/docs/configuration' },
  { label: 'Deploy Docker', to: '/docs/deployment/docker' },
];

function ProductPreview(): ReactNode {
  return (
    <div className='relative mx-auto w-full max-w-xl'>
      <div className='absolute inset-0 rounded-2xl bg-anexis-primary/20 blur-3xl' />
      <div className='relative overflow-hidden rounded-lg border border-anexis-border bg-anexis-elevated shadow-anexis'>
        <div className='flex items-center justify-between border-b border-anexis-border px-5 py-4'>
          <div className='flex items-center gap-3'>
            <img
              src='/img/anexis.png'
              alt=''
              className='h-9 w-9 object-contain'
            />
            <div>
              <p className='m-0 text-sm font-semibold text-anexis-ink'>
                Anexis Server
              </p>
              <p className='m-0 text-xs text-anexis-muted'>
                standalone · local storage
              </p>
            </div>
          </div>
          <span className='rounded-full bg-anexis-primarySoft px-3 py-1 text-xs font-semibold text-anexis-primaryStrong'>
            Healthy
          </span>
        </div>

        <div className='grid gap-0 md:grid-cols-[1fr_1.2fr]'>
          <div className='border-b border-anexis-border p-5 md:border-b-0 md:border-r'>
            <p className='mb-3 text-xs font-semibold uppercase text-anexis-muted'>
              Storage
            </p>
            <div className='space-y-3'>
              {[
                ['Objects', '128k'],
                ['Used', '42 GB'],
                ['Providers', '3'],
              ].map(([label, value]) => (
                <div
                  key={label}
                  className='flex items-center justify-between rounded-md bg-anexis-surface px-3 py-2'
                >
                  <span className='text-sm text-anexis-muted'>{label}</span>
                  <span className='text-sm font-semibold text-anexis-ink'>
                    {value}
                  </span>
                </div>
              ))}
            </div>
          </div>

          <div className='p-5'>
            <p className='mb-3 text-xs font-semibold uppercase text-anexis-muted'>
              API route
            </p>
            <pre className='m-0 overflow-x-auto rounded-md border border-anexis-border bg-anexis-code p-4 text-sm leading-7'>
              <code className='font-mono text-anexis-primary'>
                {`POST /api/files
Authorization: Bearer $ANEXIS_API_KEY
Storage-Provider: s3

201 Created`}
              </code>
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}

function FeatureCard({ label, title, description }: Feature): ReactNode {
  return (
    <article className='rounded-lg border border-anexis-border bg-anexis-elevated p-6 shadow-sm transition duration-200 hover:-translate-y-1 hover:border-anexis-primary hover:shadow-glow'>
      <p className='mb-3 text-xs font-bold uppercase tracking-wide text-anexis-primaryStrong'>
        {label}
      </p>
      <Heading as='h3' className='mb-3 text-xl font-bold text-anexis-ink'>
        {title}
      </Heading>
      <p className='m-0 text-base leading-7 text-anexis-muted'>{description}</p>
    </article>
  );
}

function HomepageHeader(): ReactNode {
  return (
    <header className='overflow-hidden bg-anexis-hero px-6 py-16 md:py-20'>
      <div className='container grid min-h-[calc(100vh-4rem)] items-center gap-12 py-8 lg:grid-cols-[0.95fr_1.05fr]'>
        <div className='max-w-2xl'>
          <div className='mb-6 inline-flex items-center gap-3 rounded-full border border-anexis-border bg-anexis-elevated px-4 py-2 text-sm font-semibold text-anexis-primaryStrong shadow-sm'>
            <span className='h-2 w-2 rounded-full bg-anexis-primary' />
            Open-source cloud file storage server
          </div>
          <Heading
            as='h1'
            className='mb-6 text-5xl font-extrabold leading-tight text-anexis-ink md:text-6xl'
          >
            Store files anywhere, expose them through one clean API.
          </Heading>
          <p className='mb-8 max-w-xl text-lg leading-8 text-anexis-muted md:text-xl'>
            Anexis Server gives self-hosters and teams a simple storage backend
            with flexible providers, migration paths and deployment modes.
          </p>
          <div className='flex flex-wrap items-center gap-4'>
            <Link
              className='button button--primary button--lg'
              to='/docs/getting-started'
            >
              Get Started
            </Link>
            <Link
              className='button button--secondary button--lg'
              to='/docs/intro'
            >
              Read the Docs
            </Link>
          </div>
        </div>

        <ProductPreview />
      </div>
    </header>
  );
}

function HomepageFeatures(): ReactNode {
  return (
    <section className='bg-anexis-band px-6 py-16 md:py-20'>
      <div className='container'>
        <div className='mb-10 flex flex-col justify-between gap-4 md:flex-row md:items-end'>
          <div>
            <p className='mb-3 text-sm font-bold uppercase tracking-wide text-anexis-primaryStrong'>
              Core capabilities
            </p>
            <Heading
              as='h2'
              className='m-0 max-w-2xl text-3xl font-bold text-anexis-ink md:text-4xl'
            >
              Practical building blocks for owning your file infrastructure.
            </Heading>
          </div>
          <Link
            to='/docs/architecture'
            className='font-semibold text-anexis-primary hover:text-anexis-primaryStrong'
          >
            Architecture overview
          </Link>
        </div>
        <div className='grid gap-5 md:grid-cols-2 lg:grid-cols-4'>
          {features.map((feature) => (
            <FeatureCard key={feature.title} {...feature} />
          ))}
        </div>
      </div>
    </section>
  );
}

function QuickStart(): ReactNode {
  return (
    <section className='bg-anexis-page px-6 py-16 md:py-20'>
      <div className='container grid gap-8 lg:grid-cols-[0.8fr_1.2fr] lg:items-start'>
        <div>
          <p className='mb-3 text-sm font-bold uppercase tracking-wide text-anexis-primaryStrong'>
            Quick start
          </p>
          <Heading
            as='h2'
            className='mb-4 text-3xl font-bold text-anexis-ink md:text-4xl'
          >
            Run it locally in one Docker command.
          </Heading>
          <p className='mb-6 text-base leading-7 text-anexis-muted'>
            Start in standalone mode, then switch providers or production
            settings when your deployment is ready.
          </p>
          <div className='flex flex-wrap gap-3'>
            {quickLinks.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                className='rounded-md border border-anexis-border bg-anexis-elevated px-4 py-2 text-sm font-semibold text-anexis-primary hover:border-anexis-primary hover:text-anexis-primaryStrong'
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>

        <pre className='m-0 overflow-x-auto rounded-lg border border-anexis-border bg-anexis-code p-6 text-sm leading-7 shadow-anexis'>
          <code className='font-mono text-anexis-primary'>
            {`docker run -d \\
  --name anexis-server \\
  -p 8080:8080 \\
  -e SERVER_MODE=standalone \\
  -e STORAGE_PROVIDER=local \\
  -e ANEXIS_API_KEY=your-key \\
  ghcr.io/gurren-software/anexis-server:latest`}
          </code>
        </pre>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title='Anexis Server - Open Source Cloud Storage'
      description='Open-source cloud file storage server built with Go. Deploy anywhere, store anything.'
    >
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <QuickStart />
      </main>
    </Layout>
  );
}
