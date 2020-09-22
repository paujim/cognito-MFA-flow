import React from 'react';
import { render } from '@testing-library/react';
import App from './App';

test('renders a app bar', () => {
  const { getByText } = render(<App />);
  const appBar = getByText(/MFA Sample/i);
  expect(appBar).toBeInTheDocument();
});
