import React from 'react';
import { render } from 'vitest-browser-react';
import App from '@/App';

test('renders learn react link', () => {
  const screen = render(<App />);
  const linkElement = screen.getByText(/learn react/i);
  expect(linkElement).toBeInTheDocument();
});
