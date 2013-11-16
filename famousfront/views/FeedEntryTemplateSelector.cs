using famousfront.viewmodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;

namespace famousfront.views
{
    class FeedEntryTemplateSelector : DataTemplateSelector
    {
        public override DataTemplate SelectTemplate(object item,    DependencyObject container)
        {
            var ws = item as FeedEntryViewModel;
            if (ws == null)
                return base.SelectTemplate(item, container);
            var fe = container as FrameworkElement;
            var mot = fe.FindResource("MediaOneTemplate") as DataTemplate;
            if (ws.HasMedia && !ws.HasDocument && !ws.HasMediaGallery)
                return mot ?? base.SelectTemplate(item, container);
            return base.SelectTemplate(item, container);
        }
    }
}
