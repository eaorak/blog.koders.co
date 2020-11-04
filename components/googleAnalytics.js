import ReactGA from "react-ga"

const GOOGLE_TRACKING_ID = 'UA-177836540-1'

export const initAnalytics = () => {
  if (!window.GA_INITIALIZED) {
    ReactGA.initialize(GOOGLE_TRACKING_ID)
    window.GA_INITIALIZED = true
  }
  ReactGA.set({ page: window.location.pathname })
  ReactGA.pageview(window.location.pathname)
}