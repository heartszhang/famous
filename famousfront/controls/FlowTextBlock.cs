using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Documents;
using System.Windows.Markup;
using System.Xml;

namespace famousfront.controls
{
    class FlowTextBlock : RichTextBox
    {
        public static readonly DependencyProperty DocumentProperty =
            DependencyProperty.Register("Document", typeof(string),
            typeof(FlowTextBlock), new FrameworkPropertyMetadata
            (null, new PropertyChangedCallback(OnDocumentChanged)));

        public new string Document
        {
            get
            {
                return (string)this.GetValue(DocumentProperty);
            }

            set
            {
                this.SetValue(DocumentProperty, value);
            }
        }

        public static void OnDocumentChanged(DependencyObject obj,
            DependencyPropertyChangedEventArgs args)
        {
            RichTextBox rtb = (RichTextBox)obj ;
            if (args.NewValue == null)
            {
              rtb.Document = new FlowDocument();
              return;
            }
            var fdoc = XamlReader.Load(new XmlTextReader(new System.IO.StringReader((string)args.NewValue))) as FlowDocument;
            var s = rtb.FindResource("FeedEntryFlowDocumentStyle") as Style;
            fdoc.Style = s;
            rtb.Document = fdoc;
            //System.Diagnostics.Debug.Write(args.NewValue);
        }
    }
}
