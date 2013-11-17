﻿using famousfront.core;
using famousfront.datamodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.viewmodels
{
  class ImageUnitViewModel : ImageBaseViewModel
  {
    internal ImageUnitViewModel(FeedMedia v)
      : base(v)
    {
    }
    public double Scale
    {
      get { return _.height > 0 ? (double)_.width / _.height : 0.0; }
    }
  }
}
