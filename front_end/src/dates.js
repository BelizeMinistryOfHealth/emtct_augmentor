import { parseISO, format } from 'date-fns';

// Convert time. date-fns does not properly convert zulu time.
// We strip the timezone information so that it can be properly parsed.
export const parseDate = (d, f) => {
  if (d) {
    // const date = parseISO(new Date(d.slice(0, -1)).toISOString());
    const date = parseISO(d);
    console.log(`date : ${date}`);
    return format(date, f);
  }
  return null;
};
