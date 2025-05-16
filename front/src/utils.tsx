const eqSet = (xs: Set<any>, ys: Set<any>) =>
    xs.size === ys.size &&
    Array.from(xs).every((x) => ys.has(x));


export {eqSet};