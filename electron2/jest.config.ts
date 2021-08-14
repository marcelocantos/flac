import type {Config} from '@jest/types';

const config: Config.InitialOptions = {
  moduleNameMapper: {
    '\\.css$': "<rootDir>/src/mock.ts",
  },
  setupFilesAfterEnv: [
    "<rootDir>/src/setupTests.ts",
  ],
  testEnvironment: 'jsdom',
  transform: {
    '\\.tsx?$': 'ts-jest',
  },
};

export default config;
